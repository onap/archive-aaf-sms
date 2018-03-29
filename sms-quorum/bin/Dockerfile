FROM alpine:3.7

LABEL name="aaf-sms"
LABEL version=1.0.0
LABEL maintainer="Girish Havaldar <hg0071052@techmahindra.com>"
#Commenting out,since no external client will invoking the Quorum client
#EXPOSE 10443

#Do not have any dependency on other modules
#RUN mkdir -p /sms/auth
ADD quorumclient /smsquorum/bin/quorumclient
RUN chmod +x /smsquorum/bin/quorumclient
#Yet to create quorumclient executable 
ENTRYPOINT ["/smsquorum/bin/quorumclient"]