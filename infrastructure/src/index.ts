import { Policy } from '@pulumi/aws/iam';
import { ComponentResource, ComponentResourceOptions, Input, Output, ResourceOptions } from '@pulumi/pulumi';

interface Args {
  /**
   * The ARN of the SES identity to allow sending on
   */
  sesIdentity: Input<string>;
}

class Mailer extends ComponentResource {
  public readonly policy: Output<string>;

  constructor(name: string, args: Args, opts?: ComponentResourceOptions) {
    super('wafflehacks:mailer:Mailer', name, { options: opts }, opts);

    const defaultResourceOptions: ResourceOptions = { parent: this };
    const { sesIdentity } = args;

    const policy = new Policy(
      `${name}-policy`,
      {
        policy: {
          Version: '2012-10-17',
          Statement: [
            {
              Effect: 'Allow',
              Action: ['ses:SendEmail', 'ses:SendRawEmail', 'ses:ListEmailIdentities'],
              Resource: [sesIdentity],
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
