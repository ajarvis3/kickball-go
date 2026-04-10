import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as path from "node:path";
import { fileURLToPath } from "node:url";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigwv2 from "aws-cdk-lib/aws-apigatewayv2";
import * as integrations from "aws-cdk-lib/aws-apigatewayv2-integrations";

export class KickballStack extends cdk.Stack {
   constructor(scope: Construct, id: string, props?: cdk.StackProps) {
      super(scope, id, props);

      // Core single-table DynamoDB design
      const table = new dynamodb.Table(this, "KickballTable", {
         tableName: "KickballTable",
         partitionKey: { name: "PK", type: dynamodb.AttributeType.STRING },
         sortKey: { name: "SK", type: dynamodb.AttributeType.STRING },
         billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
         removalPolicy: cdk.RemovalPolicy.DESTROY,
         // Enable if you want event-driven projections later
         // stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES,
      });

      //
      // GSI: At-bats by player
      //  - IndexName: GSIPlayerAtBat
      //  - Partition key: GSIPlayerAtBatPK = PLAYER#<playerId>
      //  - Sort key:      GSIPlayerAtBatSK = GAME#<gameId>#ATBAT#<seq>
      table.addGlobalSecondaryIndex({
         indexName: "GSIPlayerAtBat",
         partitionKey: {
            name: "GSIPlayerAtBatPK",
            type: dynamodb.AttributeType.STRING,
         },
         sortKey: {
            name: "GSIPlayerAtBatSK",
            type: dynamodb.AttributeType.STRING,
         },
         projectionType: dynamodb.ProjectionType.ALL,
      });

      // GSI: Games by league
      //  - IndexName: GSILeagueGame
      //  - Partition key: GSILeagueGamePK = LEAGUE#<leagueId>
      //  - Sort key:      GSILeagueGameSK = GAME#<gameId>
      table.addGlobalSecondaryIndex({
         indexName: "GSILeagueGame",
         partitionKey: {
            name: "GSILeagueGamePK",
            type: dynamodb.AttributeType.STRING,
         },
         sortKey: {
            name: "GSILeagueGameSK",
            type: dynamodb.AttributeType.STRING,
         },
         projectionType: dynamodb.ProjectionType.ALL,
      });

      // GSI: Leagues by name (supports prefix search on lowercase name)
      // - IndexName: GSILeagueByName
      // - Partition key: GSILeagueNamePK = "LEAGUE_NAME"
      // - Sort key:      GSILeagueNameSK = <lowercase league name>
      table.addGlobalSecondaryIndex({
         indexName: "GSILeagueByName",
         partitionKey: {
            name: "GSILeagueNamePK",
            type: dynamodb.AttributeType.STRING,
         },
         sortKey: {
            name: "GSILeagueNameSK",
            type: dynamodb.AttributeType.STRING,
         },
         projectionType: dynamodb.ProjectionType.ALL,
      });

      // Export table name for Lambdas / other stacks
      new cdk.CfnOutput(this, "KickballTableName", {
         value: table.tableName,
      });

      // ESM-compatible __dirname
      const __filename = fileURLToPath(import.meta.url);
      const __dirname = path.dirname(__filename);

      // Build and deploy the Go API as a Lambda function.
      // This uses Docker bundling to compile a Linux binary and place
      // the `bootstrap` executable into the asset output for a
      // custom runtime (provided.al2).
      const apiFn = new lambda.Function(this, "ApiFunction", {
         runtime: lambda.Runtime.PROVIDED_AL2,
         handler: "bootstrap",
         code: lambda.Code.fromAsset(path.join(__dirname, "../.."), {
            bundling: {
               image: cdk.DockerImage.fromRegistry("golang:1.24"),
               user: "root",
               command: [
                  "sh",
                  "-c",
                  [
                     "cd /asset-input/cmd/api",
                     "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /asset-output/bootstrap .",
                  ].join(" && "),
               ],
            },
         }),
         timeout: cdk.Duration.seconds(15),
         environment: {
            DYNAMODB_TABLE: table.tableName,
         },
      });

      // Grant the Lambda permissions to access the DynamoDB table
      table.grantReadWriteData(apiFn);

      // Create an HTTP API and integrate it with the Lambda
      const httpApi = new apigwv2.HttpApi(this, "KickballHttpApi", {
         apiName: "KickballHttpApi",
         defaultIntegration: new integrations.HttpLambdaIntegration(
            "ApiIntegration",
            apiFn,
            {
               payloadFormatVersion: apigwv2.PayloadFormatVersion.VERSION_1_0,
            },
         ),
         // Configure CORS so browsers can perform preflight requests during development
         corsPreflight: {
            allowOrigins: ["*"],
            allowMethods: [
               apigwv2.CorsHttpMethod.GET,
               apigwv2.CorsHttpMethod.POST,
               apigwv2.CorsHttpMethod.OPTIONS,
            ],
            allowHeaders: ["Content-Type", "Authorization"],
         },
         createDefaultStage: true,
      });

      new cdk.CfnOutput(this, "HttpApiUrl", {
         value: httpApi.apiEndpoint,
      });
   }
}
