FROM scratch
COPY simpleApi simpleApi
EXPOSE 5001/tcp
CMD ["./simpleApi"]