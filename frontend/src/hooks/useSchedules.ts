import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { scheduleApi } from '../services/api';
import { Schedule, ScheduleStats, StartVisitRequest, EndVisitRequest } from '../services/types';
import { useErrorStore } from '../stores/errorStore';

// Query keys
export const scheduleKeys = {
  all: ['schedules'] as const,
  lists: () => [...scheduleKeys.all, 'list'] as const,
  list: (filters: string) => [...scheduleKeys.lists(), { filters }] as const,
  details: () => [...scheduleKeys.all, 'detail'] as const,
  detail: (id: number) => [...scheduleKeys.details(), id] as const,
  stats: () => [...scheduleKeys.all, 'stats'] as const,
};

// Hooks for schedules
export const useSchedules = () => {
  return useQuery({
    queryKey: scheduleKeys.list('all'),
    queryFn: scheduleApi.getSchedules,
  });
};

export const useTodaySchedules = () => {
  return useQuery({
    queryKey: scheduleKeys.list('today'),
    queryFn: scheduleApi.getTodaySchedules,
  });
};

export const useScheduleStats = () => {
  return useQuery({
    queryKey: scheduleKeys.stats(),
    queryFn: scheduleApi.getScheduleStats,
  });
};

export const useScheduleById = (id: number) => {
  return useQuery({
    queryKey: scheduleKeys.detail(id),
    queryFn: () => scheduleApi.getScheduleById(id),
    enabled: !!id,
  });
};

// Mutations for visit management
export const useStartVisit = () => {
 const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ scheduleId, data }: { scheduleId: number; data: StartVisitRequest }) =>
      scheduleApi.startVisit(scheduleId, data),
    onSuccess: () => {
      // Invalidate and refetch schedule data
      queryClient.invalidateQueries({ queryKey: scheduleKeys.all });
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.details || error.message || 'Failed to start visit';
      useErrorStore.getState().setError(errMsg);
    },
  });
};

export const useEndVisit = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ scheduleId, data }: { scheduleId: number; data: EndVisitRequest }) =>
      scheduleApi.endVisit(scheduleId, data),
    onSuccess: () => {
      // Invalidate and refetch schedule data
      queryClient.invalidateQueries({ queryKey: scheduleKeys.all });
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.details || error.message || 'Failed to end visit';
      useErrorStore.getState().setError(errMsg);
    },
  });
};

export const useCancelVisit = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (scheduleId: number) => scheduleApi.cancelVisit(scheduleId),
    onSuccess: () => {
      // Invalidate and refetch schedule data
      queryClient.invalidateQueries({ queryKey: scheduleKeys.all });
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.details || error.message || 'Failed to cancel visit';
      useErrorStore.getState().setError(errMsg);
    },
  });
};
