## Design decision

- Original file is currently store at disk. This is because of 2 things:
    - simplicity due to time constraint
    - alignment with the CLI requirements that needs the output file path as 3rd argument
- In the future, it should be stored at object storage (s3) using pre-signed URL during conversion for scalability and robustness and only publish the file metadata if we're to expose the API to any client (frontend/mobile)