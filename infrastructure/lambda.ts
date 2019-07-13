import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as chp from "child_process";
import * as util from "util";

const chpExec = util.promisify(chp.exec);

export type LambdaOptions = {
  filepath: string;
  role: aws.iam.Role;
  handlerName: string;
  lambdaOptions?: Omit<
    aws.lambda.FunctionArgs,
    "runtime" | "code" | "timeout" | "memorySize" | "handler" | "role" | "name"
  >;
};

export const createLambdaFunction = (name, options: LambdaOptions) =>
  new aws.lambda.Function(name, {
    runtime: aws.lambda.Go1dxRuntime,
    code: new pulumi.asset.FileArchive(
      (async () => {
        await chpExec(
          `GOOS=linux GOARCH=amd64 go build -o ./dist/functions/${
            options.filepath
          }/main functions/${options.filepath}/main.go`
        );
        await chpExec(
          `zip -j ./dist/functions/${
            options.filepath
          }/main.zip ./dist/functions/${options.filepath}/main`
        );

        return `./dist/functions/${options.filepath}/main.zip`;
      })()
    ),
    timeout: 10,
    memorySize: 128,
    handler: "main",
    role: options.role.arn,
    name: options.handlerName,
    ...options.lambdaOptions
  });
