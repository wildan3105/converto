## Design decision

2) Current File Storage Approach: Local Disk
Currently, original files are stored on the local disk. This decision was made due to:

Simplicity: Given tight time constraints, storing files locally is the quickest and easiest solution to implement.
CLI Alignment: The CLI tool requires an output file path as the third argument, making local disk storage a natural fit for the current architecture.
Future Consideration: Object Storage with Pre-Signed URLs
To improve scalability and robustness, original files should eventually be stored in an object storage solution such as Amazon S3 (or a compatible alternative). The recommended approach involves:

Pre-Signed URLs: Generate pre-signed URLs during the conversion process. Clients can directly upload/download files to/from S3, reducing the server's load and improving performance.
File Metadata Exposure: Instead of returning the full file directly, the API would only expose file metadata (e.g., file name, size, download URL). This approach enhances security and flexibility, particularly when integrating with frontend or mobile clients.
Balancing Time, Product, and Quality Constraints
Time: Transitioning to object storage requires changes to the storage layer, upload/download logic, and potentially the client integration. Given current priorities, this enhancement can be deferred if the current local storage method meets short-term needs.
Product: For the CLI, the output file path requirement could still be met by using temporary local files or streaming downloads from the object storage. For API clients, exposing only metadata via pre-signed URLs provides a more modern and cloud-native experience.
Quality: Object storage improves scalability (e.g., handling large files, load balancing) and resilience (e.g., durability, redundancy). It also supports features like versioning and lifecycle management, which are valuable for production environments.
Next Steps:
Short Term: Continue using local storage while exploring edge cases and potential limitations (e.g., disk space, I/O performance).
Medium Term: Prototype object storage integration with pre-signed URLs and assess changes needed for both the CLI and API clients.
Long Term: Gradually phase out local storage, with an emphasis on a smooth transition for existing CLI users and supporting API-based uploads/downloads for new clients.


2) Known Limitation: File Upload Size (1 GB) / via API

Currently, the /api/v1/conversions endpoint can only handle file uploads up to 1 GB. This limitation arises because the API processes the entire file at once, and Fiber's BodyLimit is set as an int type, capping the maximum allowable size to 1,073,741,824 bytes (1 GB).

Future Consideration: Client-Side Chunking
To improve scalability, resilience, and support for larger files (e.g., multi-GB files), a more robust approach would be to implement client-side chunking. This method involves splitting large files into smaller chunks (ideally between 5 MB to 10 MB) on the client-side before upload. The server would then:

Accept each chunk individually.
Maintain chunk metadata to track the progress and order of received chunks.
Aggregate the chunks into a single file upon completion.
Store the final file in an object storage system (e.g., S3, MinIO).
Balancing Time, Product, and Quality Constraints
Time: Implementing client-side chunking and server-side chunk aggregation requires more development and testing time. Therefore, this enhancement should be prioritized based on product roadmap and customer needs. If the majority of use cases involve files under 1 GB, this can be a lower-priority enhancement.
Product: In the short term, communicate the 1 GB limitation to users and provide guidance on compressing or splitting files externally. For high-priority clients or scenarios requiring larger file support, chunking could be implemented as an experimental or beta feature.
Quality: By adopting a chunking approach, the system gains better fault tolerance (e.g., resumable uploads) and scalability. However, ensuring the correctness of file aggregation and handling edge cases (e.g., missing chunks, corrupted data) requires rigorous testing.
Next Steps:
Short Term: Maintain the current 1 GB limit while monitoring demand for larger uploads.
Medium Term: Research and prototype chunked file uploads. Consider using established protocols (e.g., TUS, S3 Multipart Uploads).
Long Term: Gradually roll out chunking support, starting with internal testing and expanding to a pilot group of users before a full launch.