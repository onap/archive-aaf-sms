buildclient:
	cd ../../sms-client/java && echo "Building JAVA client package" && \
	mvn package

build: buildclient
	$(MAKE) -C sms build
	$(MAKE) -C quorumclient build

deploy:
	$(MAKE) -C sms deploy
	$(MAKE) -C quorumclient deploy
