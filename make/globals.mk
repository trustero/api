SHELL := /bin/bash

TOP_TARGETS :=	all clean build deps fmt test run

PROTOS_CHANGED := $(if $(shell git diff --quiet HEAD ${REF} -- $(CURDIR)/trustero-api/protos) || echo "changed", "true",)

.PHONY: $(TOP_TARGETS) $(SUB_DIRS) _foo-build

all: clean deps build fmt

$(TOP_TARGETS): $(SUB_DIRS)

$(SUB_DIRS):
	@$(MAKE) -C $@ $(MAKECMDGOALS)


format := "%s %b %s %s\n"
MSG_LEN := $(or $(shell /usr/bin/tput cols), 120)
pad := $(shell printf '%0.1s' "."{1..$(MSG_LEN)})
EMOJI_PAD:=2
# @arg $1: make target
# @arg $2: module
# @arg $3: message
define outmsg
	$(eval is_build := $(shell if [[ "$@" == *"build"* ]]; then echo "true"; fi))
	$(eval is_release := $(shell if [[ "$@" == *"release"* ]]; then echo "true"; fi))
	$(eval is_gen := $(shell if [[ "$@" == *"gen"* ]]; then echo "true"; fi))
	$(eval is_clean := $(shell if [[ "$@" == *"clean"* ]]; then echo "true"; fi))
	$(eval is_deps := $(shell if [[ "$@" == *"deps"* ]]; then echo "true"; fi))
	$(eval is_test := $(shell if [[ "$@" == *"test"* ]]; then echo "true"; fi))
	$(eval module := $(shell echo "[$(strip $1)]") )

	$(eval emoji := $(if $(is_deps),					"⬇️ ", \
                        $(if $(is_build),				"⚙️ ", \
                        $(if $(is_clean),				"🧹",	\
                        $(if $(filter "$@", "fmt" ),	 	"✍️ ",	\
                        $(if $(filter "$@", "lint" ),	"🔎",	\
                        $(if $(is_test),				"🔬",	\
                        $(if $(filter "$@", "run" ),		"🚀",	\
                        $(if $(is_release),				"🐿 ",	\
                        $(if $(is_gen),	 				"🦄",	\
                        								"$@"
                        								))))))))))
	$(eval msg := $(if $2, $2, \
						$(if $(is_deps),				"Fetching dependencies", \
						$(if $(is_build), 				"Building module", \
						$(if $(is_clean), 				"Cleaning build files", \
						$(if $(filter "$@", "fmt" ), 	"Applying code style", \
                        $(if $(filter "$@", "lint" ),	"Linting",	\
						$(if $(is_test),			 	"Testing", \
						$(if $(filter "$@", "run" ), 	"Running", \
						$(if $(is_release), 			"Releasing", \
						$(if $(is_gen),					"Generating code", \
														"$@" )))))))))))

	$(eval extra_padding := $(if $(filter "$@", "deps" ), $(EMOJI_PAD), \
							$(if $(is_build), $(EMOJI_PAD), \
							$(if $(filter "$@", "fmt" ),	$(EMOJI_PAD), 0))))


	$(eval msg_no_pad := $(shell printf $(format) $(emoji) $(msg) "" $(module)))
	$(eval msg_no_pad_len := $(shell echo $(msg_no_pad) | wc -c))
	$(eval pad_len := $(shell echo $$(( $(MSG_LEN) - $(msg_no_pad_len) +$(extra_padding)))))
	@printf $(format) \
		$(emoji) \
		"\033[1m"$(msg)"\033[0m" \
		$(shell printf "%*.*s" 0 $(pad_len) $(pad)) \
		$(module)

endef

define not-supported-msg
	$(eval msg := $(shell printf $(format) "🚧" "Nothing to be done for '$@'" "[$(strip $1)]"))
	$(eval msg_len := $(shell echo $(msg) | wc -c))
	$(eval pad_len := $(shell echo $$(( $(MSG_LEN) - $(msg_len) - $(EMOJI_PAD)))))
	@printf $(format) "🚧" \
			"\033[1mNothing to be done for '$@'\033[0m" \
			$(shell printf "%*.*s" 0 $(pad_len) $(pad)) \
			"[$(strip $1)]"
endef
