FROM ghcr.io/charmbracelet/vhs

RUN apt-get update && \
    apt-get -y install vim

WORKDIR /vhs

COPY tapes/. /tapes
COPY notes/. /notes

# Copy the built binary and documentation files from the builder stage
COPY --from=note /app/note /bin/note

RUN export EDITOR=/bin/vim && \
    for file in /tapes/*.tape; do \
       vhs $file; \
    done
