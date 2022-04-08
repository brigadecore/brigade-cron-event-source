load('ext://min_k8s_version', 'min_k8s_version')
min_k8s_version('1.18.0')

trigger_mode(TRIGGER_MODE_MANUAL)

load('ext://namespace', 'namespace_create')
namespace_create('brigade-cron-event-source')
k8s_resource(
  new_name = 'namespace',
  objects = ['brigade-cron-event-source:namespace'],
  labels = ['brigade-cron-event-source']
)

k8s_resource(
  new_name = 'event-source',
  objects = ['brigade-cron-event-source:secret'],
  labels = ['brigade-cron-event-source']
)

docker_build(
  'brigadecore/brigade-cron-event-source', '.',
  only = [
    'config.go',
    'go.mod',
    'go.sum',
    'main.go',
    'validation.go'
  ],
  match_in_env_vars = True
)

k8s_yaml(
  helm(
    './charts/brigade-cron-event-source',
    name = 'brigade-cron-event-source',
    namespace = 'brigade-cron-event-source',
    set = [
      'brigade.apiToken=' + os.environ['BRIGADE_API_TOKEN']
    ]
  )
)
