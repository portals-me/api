import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

export const createLambdaSubscription = async (
  name: string,
  options: {
    function: aws.lambda.Function;
    snsTopicArn: string;
  }
) => {
  const subscription = new aws.sns.TopicSubscription(name, {
    protocol: "lambda",
    endpoint: options.function.arn,
    topic: options.snsTopicArn as any
  });

  const permission = new aws.lambda.Permission(name, {
    function: options.function.name,
    action: "lambda:InvokeFunction",
    principal: "sns.amazonaws.com"
  });

  return subscription;
};
