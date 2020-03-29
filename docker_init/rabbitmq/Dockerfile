FROM rabbitmq:3.8.3-management

ADD plugins/rabbitmq_delayed_message_exchange-3.8.0.ez /opt/rabbitmq/plugins/

RUN rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange