## Integration Testing Approach  

The integration tests are designed to validate the core functionalities of the application by simulating real user-initiated requests. Instead of using extensive mocking or stubbing for external dependencies, these tests interact directly with the running application, ensuring that the entire request-handling pipeline—from API endpoints to database interactions—is functioning as expected.  

### Rationale Behind This Approach  

Given the time constraints, setting up a fully isolated test environment with mocked dependencies (e.g., mocking the database, RabbitMQ, or file storage interactions) would require additional effort for configuration and maintenance. While such an approach is beneficial for unit tests, the priority here is to ensure that the end-to-end behavior of the system works correctly without introducing unnecessary complexity.  

By executing actual HTTP requests against the running service, this approach:  

- **Reduces Manual Testing Effort** – Automates the validation of essential workflows, reducing the need for repetitive manual testing.  
- **Validates the Real System Behavior** – Since these tests use real API calls, they provide confidence that all integrated components (request handling, database transactions, file processing, etc.) work as expected in a real-world scenario.  
- **Ensures Data Flow Integrity** – By interacting with an actual database (MongoDB) and processing real HTTP requests, the tests verify that data is correctly stored, retrieved, and updated throughout the application's lifecycle.  
- **Minimizes Setup Overhead** – Mocking various dependencies (e.g., MongoDB, RabbitMQ, S3) can be time-consuming, and configuring a test environment to mimic production conditions would add extra complexity. This approach provides a practical balance between thorough testing and fast implementation.  
- **Focuses on the Happy Flow** – Given the limited time, the primary focus is on ensuring that the core functionalities (file upload, conversion process, retrieval of results) operate correctly under normal conditions. This ensures that the service is stable before considering edge cases and failure scenarios.  

### Trade-offs and Future Improvements  

While this approach effectively validates end-to-end functionality, it does not fully cover edge cases, failure scenarios, or performance considerations. In the future, additional test layers (such as unit tests with mocks and integration tests with controlled failure conditions) can be introduced to enhance test coverage and robustness.  
