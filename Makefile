PROJECT = file-organizer

.PHONY: all build clean install uninstall help

# Default target
all: clean install

# Build the project
build:
	@if [ ! -f ${PROJECT} ]; then printf "Building ${PROJECT}..." \
	&& go build -o ${PROJECT} && printf "Done\n"; fi

# Clean the built files
clean:
	@if [ -f ${PROJECT} ]; then printf "Cleaning ${PROJECT}..." \
	&& rm -f ${PROJECT} && printf "Done\n"; fi

# Install the project
install:
	@make build
	@printf "Installing ${PROJECT}..."
	@mkdir -p ~/.config/${PROJECT}/
	@cp config.json ~/.config/${PROJECT}/
	@mv ${PROJECT} ~/.local/bin/
	@printf "Done\n"

# Uninstall the project
uninstall:
	@printf "Uninstalling ${PROJECT}..."
	@rm -rf ~/.config/${PROJECT}/
	@rm -f ~/.local/bin/${PROJECT}
	@printf "Done\n"

# Display help screen
help:
	@printf "Usage: make [target]\n\n"
	@printf "Targets:\n"
	@printf "  all        : Clean and install the project\n"
	@printf "  build      : Build the project executable\n"
	@printf "  clean      : Remove the built executable\n"
	@printf "  install    : Build and install the project\n"
	@printf "  uninstall  : Uninstall the project\n"
	@printf "  help       : Show this help message\n"
