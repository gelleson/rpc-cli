FROM scratch

COPY rpc-cli /rpc-cli

ENTRYPOINT ["/rpc-cli"]
