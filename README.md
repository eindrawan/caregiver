# Caregivers Shift Tracker App

## Setup Instructions

To run the app locally after pulling/cloning this repository:

### Automatic Setup (Recommended)
1. Ensure you have Go (1.21+) and Node.js (18+) installed.
2. Run the `setup.sh` script in the root directory (works on Linux, macOS, Git Bash, or WSL on Windows):
   ```
   chmod +x setup.sh
   ./setup.sh
   ```
   This will:
   - Install backend dependencies (`go mod download`).
   - Build the backend executable.
   - Start the backend server in the background.
   - Install frontend dependencies (`npm install`).
   - Start the frontend development server with Expo.

   The backend will run on `http://localhost:8080`, and the frontend will open in your browser or Expo Go app for mobile testing.

   For Windows cmd.exe, use manual setup or run via Git Bash/WSL.

### Manual Setup
#### Backend (Go)
1. Navigate to the backend directory:
   ```
   cd backend
   ```
2. Install dependencies:
   ```
   go mod download
   ```
3. Build the server:
   ```
   go build -o bin/server ./cmd/server
   ```
4. Run the server:
   ```
   ./bin/server
   ```
   The API will be available at `http://localhost:8080`.

#### Frontend (React Native with Expo)
1. Navigate to the frontend directory:
   ```
   cd frontend
   ```
2. Install dependencies:
   ```
   npm install
   ```
3. Start the development server:
   ```
   npx expo start
   ```
   - For web: Open in browser.
   - For mobile: Scan QR code with Expo Go app.

#### Environment Configuration
- Copy `frontend/.env.example` to `frontend/.env` and update any API endpoints or secrets if needed (e.g., backend URL).
- The backend uses SQLite (`caregiver_shift_tracker.db`) in the backend directory; no additional setup required for local development.

### Docker Setup (Optional)
If Docker is installed, use `docker-compose up` to run both services in containers. See `docker-compose.yml` for details.

## Tech Stack and Key Decisions

### Frontend
- **React Native with Expo**: Enables cross-platform development for web, iOS, and Android from a single codebase. Expo simplifies setup, handles native dependencies, and provides tools like over-the-air updates.
- **TypeScript**: Ensures type safety, reduces runtime errors, and improves developer experience in a large codebase.
- **State Management**: React Query for server state (caching, mutations), Zustand/Redux Toolkit for global state if needed.
- **UI Components**: Custom atomic design pattern (atoms, molecules, organisms) for reusability. Styling follows a modern, teal/orange theme with Roboto font.
- **Key Decision**: Expo chosen for rapid prototyping and universal app support without ejecting for native code.

### Backend
- **Go (Golang) with Gin Framework**: Go provides high performance, concurrency, and simplicity for API services. Gin is lightweight and fast for routing/middleware.
- **Database**: SQLite for local development and simplicity; easy to swap for PostgreSQL/MySQL in production.
- **Structure**: Clean separation with handlers, services, repositories, and models following MVC-like patterns.
- **API**: RESTful endpoints with Swagger documentation (`/docs/swagger.json`). JWT auth ready for extension.
- **Key Decision**: Go selected for backend efficiency and low resource usage, ideal for mobile app APIs. SQLite keeps setup minimal.

### Overall Architecture
- Monorepo with separate `backend/` and `frontend/` directories.
- Communication: Frontend calls backend API via Axios/fetch.
- Testing: Jest for frontend, Go's built-in testing for backend.
- Deployment: Dockerfiles and docker-compose.yml provided for containerization.

This stack balances development speed, maintainability, and scalability.

## Assumptions Made
- User authentication and authorization are not implemented.
- Users have GPS enabled device, even tough i prepare some fallback
- I feel the UI design is too much variations for a simple app, so I made justification so that the UI component design is more appropriate
- Setup script is bash-compatible for cross-platform use (Linux/macOS/Git Bash/WSL); Windows cmd users should use manual instructions.

## Optional Notes: Improvements and Future Additions
If given more time, I would do:
- [ ] Checking the geolocation compared to the client's address and alerting the user if the distance is too far
- [ ] Get the adress of the geolocation if its too far from client's address
- [ ] Photo capture for visit verification
- [ ] Offline data synchronization

This provides a solid foundation for the Caregivers shift tracking app. Feel free to extend as needed!