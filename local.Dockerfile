FROM scratch
COPY dist/cnfuzz /
EXPOSE 8080
ENTRYPOINT ["/cnfuzz"]
