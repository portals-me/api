import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as lambda from "./lambda";

export const hookDynamoDB = (
  name: string,
  options: {
    topic: aws.sns.TopicArgs;
    lambda: (topicArn: pulumi.Output<string>) => lambda.LambdaOptions;
    table: aws.dynamodb.Table;
  }
) => {
  const topic = new aws.sns.Topic(name, options.topic);

  const hookFunction = lambda.createLambdaFunction(
    name,
    options.lambda(topic.arn),
    {
      dependsOn: [topic]
    }
  );

  const subscription = new aws.dynamodb.TableEventSubscription(
    name,
    options.table,
    hookFunction,
    {
      startingPosition: "TRIM_HORIZON"
    },
    {
      dependsOn: [topic, hookFunction, options.table]
    }
  );

  return {
    topic,
    hookFunction,
    subscription
  };
};
