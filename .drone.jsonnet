local pipeline = import 'pipeline.libsonnet';
local name = 'autoscaler';

[
  pipeline.test('linux', 'amd64'),
  pipeline.build(name, 'linux', 'amd64'),
  pipeline.build(name, 'linux', 'arm64'),
  pipeline.build(name, 'linux', 'arm'),
  pipeline.notifications(depends_on=[
    'linux-amd64',
    'linux-arm64',
    'linux-arm',
  ]),
]
