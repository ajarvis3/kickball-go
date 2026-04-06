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
    // GSI2 — At-bats by player
    //
    //   GSI2PK = PLAYER#<playerId>
    //   GSI2SK = GAME#<gameId>#ATBAT#<seq>
    //
    // This is the only GSI you need right now.
    //
    table.addGlobalSecondaryIndex({
      indexName: "GSI2",
      partitionKey: { name: "GSI2PK", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "GSI2SK", type: dynamodb.AttributeType.STRING },
      projectionType: dynamodb.ProjectionType.ALL,
    });

    // Export table name for Lambdas / other stacks
    new cdk.CfnOutput(this, "KickballTableName", {
      value: table.tableName,
    });
  }
}