# Distroless publishes no versioned tags
# hadolint ignore=DL3006
FROM gcr.io/distroless/static
COPY vulcain /
CMD ["/vulcain"]
EXPOSE 80 443

