FROM alpine:3.7

LABEL name="aaf-sms"
LABEL version=1.0.0
LABEL maintainer="Girish Havaldar <hg0071052@techmahindra.com>"

ADD quorumclient /smsquorum/bin/quorumclient
RUN chmod +x /smsquorum/bin/quorumclient
ENTRYPOINT ["/smsquorum/bin/quorumclient"]
