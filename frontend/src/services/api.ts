import axios from 'axios';
import {
  Schedule,
  ScheduleStats,
  Visit,
  Task,
  Client,
  ApiResponse,
  StartVisitRequest,
  EndVisitRequest,
  UpdateTaskRequest
} from './types';

// Configure base URL - adjust this to match your backend
const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://192.168.68.103:8080/api/v1';

// Mock caregiver ID - in real app this would come from authentication
const CURRENT_CAREGIVER_ID = 1;

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
apiClient.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API Response Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export const scheduleApi = {
  // Get all schedules
  getSchedules: async (): Promise<Schedule[]> => {
    const response = await apiClient.get<ApiResponse<Schedule[]>>('/schedules');
    return response.data.data;
  },

  // Get today's schedules
  getTodaySchedules: async (): Promise<Schedule[]> => {
    const response = await apiClient.get<ApiResponse<Schedule[]>>(`/schedules/today?caregiver_id=${CURRENT_CAREGIVER_ID}`);
    return response.data.data;
  },

  // Get schedule by ID
  getScheduleById: async (id: number): Promise<Schedule> => {
    const response = await apiClient.get<ApiResponse<Schedule>>(`/schedules/${id}`);
    return response.data.data;
  },

  // Get schedule statistics
  getScheduleStats: async (): Promise<ScheduleStats> => {
    const response = await apiClient.get<ApiResponse<ScheduleStats>>(`/schedules/stats?caregiver_id=${CURRENT_CAREGIVER_ID}`);
    return response.data.data;
  },



  // Start a visit
  startVisit: async (scheduleId: number, data: StartVisitRequest): Promise<Visit> => {
    const response = await apiClient.post<ApiResponse<Visit>>(`/schedules/${scheduleId}/start`, data);
    return response.data.data;
  },

  // End a visit
  endVisit: async (scheduleId: number, data: EndVisitRequest): Promise<Visit> => {
    const response = await apiClient.post<ApiResponse<Visit>>(`/schedules/${scheduleId}/end`, data);
    return response.data.data;
  },

  // Cancel a visit
  cancelVisit: async (scheduleId: number): Promise<void> => {
    await apiClient.post<ApiResponse<void>>(`/schedules/${scheduleId}/cancel`);
  },
};

export const taskApi = {
  // Get task by ID
  getTaskById: async (id: number): Promise<Task> => {
    const response = await apiClient.get<ApiResponse<Task>>(`/tasks/${id}`);
    return response.data.data;
  },

  // Update task status
  updateTaskStatus: async (id: number, data: UpdateTaskRequest): Promise<Task> => {
    const response = await apiClient.put<ApiResponse<Task>>(`/tasks/${id}`, data);
    return response.data.data;
  },
};

export const visitApi = {
  // Get visit by schedule ID
  getVisitByScheduleId: async (scheduleId: number): Promise<Visit> => {
    const response = await apiClient.get<ApiResponse<Visit>>(`/visits/schedule/${scheduleId}`);
    return response.data.data;
  },
};

export const clientApi = {
  // Get all clients
  getClients: async (): Promise<Client[]> => {
    const response = await apiClient.get<ApiResponse<{ clients: Client[]; count: number }>>('/clients');
    return response.data.data.clients;
  },

  // Get client by ID
  getClientById: async (id: number): Promise<Client> => {
    const response = await apiClient.get<ApiResponse<{ client: Client }>>(`/clients/${id}`);
    return response.data.data.client;
  },

  // Search clients
  searchClients: async (query: string): Promise<Client[]> => {
    const response = await apiClient.get<ApiResponse<{ clients: Client[]; count: number; query: string }>>(`/clients/search?q=${encodeURIComponent(query)}`);
    return response.data.data.clients;
  },

  // Create client
  createClient: async (client: Omit<Client, 'id' | 'created_at' | 'updated_at'>): Promise<Client> => {
    const response = await apiClient.post<ApiResponse<{ client: Client }>>('/clients', client);
    return response.data.data.client;
  },

  // Update client
  updateClient: async (id: number, client: Partial<Omit<Client, 'id' | 'created_at' | 'updated_at'>>): Promise<Client> => {
    const response = await apiClient.put<ApiResponse<{ client: Client }>>(`/clients/${id}`, client);
    return response.data.data.client;
  },

  // Delete client
  deleteClient: async (id: number): Promise<void> => {
    await apiClient.delete(`/clients/${id}`);
  },
};
