
PROJECT = file-organizer

.PHONY: all build clean install uninstall

all: clean install

build:
	@if [ ! -f ${PROJECT} ]; then printf "Building ${PROJECT}..." \
	&& go build -o ${PROJECT} && printf "Done\n"; fi

clean:
	@if [ -f ${PROJECT} ]; then printf "Cleaning ${PROJECT}..." \
	&& rm -f ${PROJECT} && printf "Done\n"; fi

install:
	@make build
	@printf "Installing ${PROJECT}..."
	@mkdir -p ~/.config/${PROJECT}/
	@cp config.json ~/.config/${PROJECT}/
	@mv ${PROJECT} ~/.local/bin/
	@printf "Done\n"

uninstall:
	@printf "Uninstalling ${PROJECT}..."
	@rm -rf ~/.config/${PROJECT}/
	@rm -f ~/.local/bin/${PROJECT}
	@printf "Done\n"