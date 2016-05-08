IMAGE=sejvlond/go-ultimate-server

go-ultimate-server:
	docker run \
    	--rm \
	    -v `pwd`:/src \
	    sejvlond/go-ultimate-server_build \
	    github.com/sejvlond/go-ultimate-server

build: go-ultimate-server
	docker build -t ${IMAGE} .

push: build
	docker push ${IMAGE}

run:
	docker run \
		--rm \
		-it \
		-v `pwd`/conf:/www/ultimate-server/conf \
		-P \
		${IMAGE}

clean:
	sudo rm go-ultimate-server

