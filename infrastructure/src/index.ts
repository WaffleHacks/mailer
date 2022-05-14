import { Policy } from '@pulumi/aws/iam';
import {
  ComponentResource,
  ComponentResourceOptions,
  Input,
  Output,
  ResourceOptions,
  interpolate,
} from '@pulumi/pulumi';

interface Args {
  /**
   * The domains emails can be sent from
   */
  fromDomains: Input<string>[];
}

class Mailer extends ComponentResource {
  public readonly policy: Output<string>;

  constructor(name: string, args: Args, opts?: ComponentResourceOptions) {
    super('wafflehacks:mailer:Mailer', name, { options: opts }, opts);

    const defaultResourceOptions: ResourceOptions = { parent: this };
    const { fromDomains } = args;

    const policy = new Policy(
      `${name}-policy`,
      {
        policy: {
          Version: '2012-10-17',
          Statement: [
            {
              Effect: 'Allow',
              Action: ['ses:SendEmail', 'ses:SendRawEmail'],
              Resource: '*',
              Condition: {
                'ForAllValues:StringLike': {
                  'ses:FromAddress': fromDomains.map((d) => interpolate`*@${d}`),
                },
              },
            },
            {
              Effect: 'Allow',
              Action: ['ses:ListEmailIdentities'],
              Resource: '*',
            },
          ],
        },
      },
      defaultResourceOptions,
    );

    this.policy = policy.name;
    this.registerOutputs();
  }
}

export default Mailer;
