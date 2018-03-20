FROM alpine:3.7
MAINTAINER aaf-sms

LABEL name="aaf-sms"
LABEL version=1.0.0
LABEL maintainer="vamshi krishna <vn00480215@techmahindra.com>"

EXPOSE 10443

ADD auth /sms/auth/
ADD sms /sms/bin/sms
RUN chmod +x /sms/bin/sms

ENTRYPOINT ["/sms/bin/sms"]