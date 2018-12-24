# Inspired by https://github.com/minds-ai/zoom-drive-connector/blob/master/Makefile
NAME 	= "thesaurize-bot"
TAG		= $(shell git --no-pager log -1 --pretty=%H)
IMG		= ${NAME}:${TAG}

build:
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${NAME}:latest