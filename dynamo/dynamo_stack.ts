import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";

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

      // Export table name for Lambdas / other stacks
      new cdk.CfnOutput(this, "KickballTableName", {
         value: table.tableName,
      });
   }
}
