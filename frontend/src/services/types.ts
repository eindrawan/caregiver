export interface Client {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  address: string;
  city: string;
  state: string;
  zip_code: string;
  latitude: number;
  longitude: number;
  notes?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Schedule {
  id: number;
  client_id: number;
  service_name?: string;
  caregiver_id: number;
  start_time: string;
  end_time: string;
  status: 'scheduled' | 'in_progress' | 'completed' | 'missed' | 'cancelled';
  notes: string;
  created_at: string;
  updated_at: string;
  client?: Client;
  visit?: Visit;
  tasks?: Task[];
}

export interface Visit {
  id: number;
  schedule_id: number;
  start_time?: string;
  end_time?: string;
  start_latitude?: number;
  start_longitude?: number;
  end_latitude?: number;
  end_longitude?: number;
  status: 'not_started' | 'in_progress' | 'completed';
  notes?: string;
}

export interface Task {
  id: number;
  schedule_id: number;
  name: string;
  description: string;
  status: 'pending' | 'completed' | 'not_completed';
  reason?: string;
}

export interface ScheduleStats {
  total: number;
  missed: number;
  upcoming: number;
  completed: number;
}

export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
  error?: string;
}

export interface LocationData {
  latitude: number;
  longitude: number;
}

export interface StartVisitRequest {
  start_latitude: number;
  start_longitude: number;
}

export interface EndVisitRequest {
  end_latitude: number;
  end_longitude: number;
  notes?: string;
}

export interface UpdateTaskRequest {
  status: 'completed' | 'not_completed';
  reason?: string;
}
