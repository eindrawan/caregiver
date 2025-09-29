import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { taskApi } from '../services/api';
import { Task, UpdateTaskRequest } from '../services/types';
import { useErrorStore } from '../stores/errorStore';

// Query keys for tasks
export const taskKeys = {
  all: ['tasks'] as const,
  lists: () => [...taskKeys.all, 'list'] as const,
  list: (filters: string) => [...taskKeys.lists(), { filters }] as const,
  details: () => [...taskKeys.all, 'detail'] as const,
  detail: (id: number) => [...taskKeys.details(), id] as const,
};

// Hook to get a single task by ID
export const useTaskById = (id: number) => {
  return useQuery({
    queryKey: taskKeys.detail(id),
    queryFn: () => taskApi.getTaskById(id),
    enabled: !!id,
  });
};

// Hook to update task status
export const useUpdateTaskStatus = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ taskId, data }: { taskId: number; data: UpdateTaskRequest }) =>
      taskApi.updateTaskStatus(taskId, data),
    onSuccess: (updatedTask) => {
      // Update the specific task in cache
      queryClient.setQueryData(taskKeys.detail(updatedTask.id), updatedTask);
      
      // Invalidate schedule queries to refresh task lists
      queryClient.invalidateQueries({ queryKey: ['schedules'] });
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.details || error.message || 'Failed to update task status';
      useErrorStore.getState().setError(errMsg);
    },
  });
};

// Hook to update multiple tasks (for batch operations)
export const useUpdateMultipleTasks = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (updates: { taskId: number; data: UpdateTaskRequest }[]) => {
      const promises = updates.map(({ taskId, data }) => 
        taskApi.updateTaskStatus(taskId, data)
      );
      return Promise.all(promises);
    },
    onSuccess: (updatedTasks) => {
      // Update each task in cache
      updatedTasks.forEach((task) => {
        queryClient.setQueryData(taskKeys.detail(task.id), task);
      });
      
      // Invalidate schedule queries to refresh task lists
      queryClient.invalidateQueries({ queryKey: ['schedules'] });
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.details || error.message || 'Failed to update multiple tasks';
      useErrorStore.getState().setError(errMsg);
    },
  });
};
