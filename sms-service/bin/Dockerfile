FROM alpine:3.7
MAINTAINER onap-aaf-sms

ENV VAULT_VERSION 0.9.4
ENV CONSUL_VERSION 1.0.1

LABEL name="sms-vault-consul"
LABEL version=0.0.1
LABEL maintainer="vamshi krishna <vn00480215@techmahindra.com>"


 ADD https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip /tmp/vault.zip
 RUN unzip -d /bin /tmp/vault.zip && \
 chmod 755 /bin/vault && \
 rm /tmp/vault.zip

 ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip /tmp/consul.zip
 RUN unzip -d /bin /tmp/consul.zip && \
 chmod 755 /bin/consul && \
 rm /tmp/consul.zip


EXPOSE 8200
VOLUME "/config"
ADD ./config /config/

ENTRYPOINT ["/bin/consul"]
CMD ["agent", "-config-file", "config/consul.json"]


ENTRYPOINT ["/bin/vault"]
CMD ["server", "-config", "config/vault.hcl"]

