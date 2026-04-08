import * as cdk from 'aws-cdk-lib';
import { KickballStack } from '../dynamo_stack.js';

const app = new cdk.App();
new KickballStack(app, 'KickballStack');
