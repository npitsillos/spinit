*Pull requests, bug reports, and all other forms of contribution are welcomed and highly encouraged!* :octocat:

# Development Flow
## Building
Fork and clone the repository locally. To build `spinit` simply run `make build`, which will build a binary in the `./bin/` directory. Alternatively you can run `go build -o <exe-name>`.

## Testing
There are two different test types for `spinit`, unit tests and integration tests. Every feature pr code change should be tested where applicable. To run unit or integration tests simply run `make unit-test` or `make integration-tests`.

# Pull Requests
## Commit Message Guidelines

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification to ensure clear and descriptive commit messages. Following these guidelines helps maintain a consistent commit history, making it easier to track changes and understand the project’s evolution.

### Format

Please use the following structure for your commit messages:
- **type**: Describes the purpose of the change. Common types include:
  - `feat`: A new feature
  - `fix`: A bug fix
  - `docs`: Documentation updates
  - `style`: Code style changes (formatting, etc.)
  - `refactor`: Code changes that don’t affect functionality
  - `test`: Adding or updating tests
  - `chore`: Other changes, like build or tooling updates

- **scope**: Specifies what part of the codebase is affected (optional but recommended, e.g., `auth`, `ui`, `api`).

- **short description**: A brief summary of the changes, written in imperative mood (e.g., "add login button").

### Examples

- `feat(auth): add new OAuth login flow`
- `fix(ui): correct button alignment on mobile`
- `docs: update README with setup instructions`
- `refactor(database): improve query performance`

### Additional Notes

For more substantial or complex changes, consider adding more detail in the commit body. Separate the body from the header with a blank line. Use the body to explain *why* the change was made, which can be helpful during code reviews and debugging.

Following these guidelines will make your contributions easier to integrate and understand. Thank you for your effort!