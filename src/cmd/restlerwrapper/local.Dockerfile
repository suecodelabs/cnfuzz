FROM scratch
COPY dist/restlerwrapper /
ENTRYPOINT ["/restlerwrapper"]
