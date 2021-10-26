EXE = podbit

UISRC    = ui/ui.go ui/input.go colors/colors.go ui/library.go ui/player.go ui/queue.go ui/download.go ui/tray.go
UICOMPS  = ui/components/menu.go ui/components/table.go
SOUNDSRC = sound/sound.go sound/queue.go
DATASRC  = data/data.go data/queue.go data/db.go data/cache.go
SRC = main.go ver.go ${INPUTSRC} ${UISRC} ${DATASRC} ${UICOMPS} ${SOUNDSRC}

ifndef PREFIX
	PREFIX = /usr/local
endif

${EXE}: ${SRC}
	CGO_LDFLAGS_ALLOW=".*" go build

check:
	CGO_LDFLAGS_ALLOW=".*" go run -race . 2>race.log
clean:
	go clean

install: ${EXE}
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f ${EXE} ${DESTDIR}${PREFIX}/bin/
	chmod 755 ${DESTDIR}${PREFIX}/bin/${EXE}

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/${EXE}

.PHONY = check clean install uninstall
