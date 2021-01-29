import { App, Construct, Duration, RemovalPolicy, Stack, StackProps } from '@aws-cdk/core';
import { AttributeType, BillingMode, ProjectionType, StreamViewType, Table } from '@aws-cdk/aws-dynamodb';
import { Code, Function, Runtime, StartingPosition, Tracing } from '@aws-cdk/aws-lambda';
import { DynamoEventSource } from '@aws-cdk/aws-lambda-event-sources';
import { Rule, RuleTargetInput, Schedule } from '@aws-cdk/aws-events';
import { LambdaFunction } from '@aws-cdk/aws-events-targets';
import { HttpApi, PayloadFormatVersion } from '@aws-cdk/aws-apigatewayv2';
import { LambdaProxyIntegration } from '@aws-cdk/aws-apigatewayv2-integrations';

interface SchedulerProps extends StackProps {
  name: string;
  version: string;
}

class SchedulerStack extends Stack {
  constructor(scope: Construct, id: string, props: SchedulerProps) {
    super(scope, id, props);

    const schedulerTable = new Table(this, 'Table', {
      tableName: `${props.name}-${props.version}`,
      removalPolicy: RemovalPolicy.RETAIN,
      billingMode: BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      partitionKey: {
        name: 'id',
        type: AttributeType.STRING
      },
      stream: StreamViewType.OLD_IMAGE
    });

    schedulerTable.addGlobalSecondaryIndex({
      indexName: 'ix_due_at',
      partitionKey: {
        name: 'dummy',
        type: AttributeType.STRING
      },
      sortKey: {
        name: 'dueAt',
        type: AttributeType.NUMBER
      },
      projectionType: ProjectionType.KEYS_ONLY
    });

    const graphqlLambda = new Function(this, 'GraphQLFunction', {
      functionName: `${props.name}-graphql-${props.version}`,
      handler: 'main',
      runtime: Runtime.GO_1_X,
      memorySize: 3008,
      timeout: Duration.seconds(30),
      tracing: Tracing.ACTIVE,
      code: Code.fromAsset(`./../graphql/dist`),
      environment: {
        SCHEDULER_TABLE_NAME: schedulerTable.tableName
      }
    });

    schedulerTable.grantReadWriteData(graphqlLambda);

    const collectorLambda = new Function(this, 'CollectorFunction', {
      functionName: `${props.name}-collector-${props.version}`,
      handler: 'main',
      runtime: Runtime.GO_1_X,
      memorySize: 3008,
      timeout: Duration.minutes(15),
      tracing: Tracing.ACTIVE,
      code: Code.fromAsset(`./../collector/dist`),
      environment: {
        SCHEDULER_TABLE_NAME: schedulerTable.tableName
      }
    });

    schedulerTable.grantReadWriteData(collectorLambda);

    new Rule(this, 'SchedulerRule', {
      schedule: Schedule.rate(Duration.minutes(1)),
      targets: [new LambdaFunction(collectorLambda, { event: RuleTargetInput.fromObject({ })})]
    });

    const workerLambda = new Function(this, 'WorkerFunction', {
      functionName: `${props.name}-worker-${props.version}`,
      handler: 'main',
      runtime: Runtime.GO_1_X,
      memorySize: 3008,
      timeout: Duration.minutes(15),
      tracing: Tracing.ACTIVE,
      code: Code.fromAsset(`./../worker/dist`)
    });

    workerLambda.addEventSource(new DynamoEventSource(schedulerTable, {
      startingPosition: StartingPosition.LATEST
    }));

    schedulerTable.grantStreamRead(workerLambda);

    const integration = new LambdaProxyIntegration({
      handler: graphqlLambda,
      payloadFormatVersion: PayloadFormatVersion.VERSION_2_0
    });

    new HttpApi(this, 'Api', {
      apiName: `${props.name}-api-${props.version}`,
      defaultIntegration: integration
    });
  }
}

const app = new App();

new SchedulerStack(app, 'scheduler-v1', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION
  },
  name: 'scheduler',
  version: 'v1'
});

app.synth();
