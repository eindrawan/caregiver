// Utility functions

export const formatTime = (timeString: string): string => {
  return new Date(timeString).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  });
};

export const formatDate = (timeString: string): string => {
  return new Date(timeString).toLocaleDateString('en-US', {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  });
};

export const formatDateTime = (timeString: string): string => {
  const date = new Date(timeString);
  return `${formatDate(timeString)} ${formatTime(timeString)}`;
};

export const getStatusColor = (status: string): string => {
  switch (status) {
    case 'completed':
      return '#10B981'; // success
    case 'in_progress':
      return '#F97316'; // accent
    case 'missed':
    case 'cancelled':
      return '#EF4444'; // error
    case 'scheduled':
    default:
      return '#6B7280'; // gray
  }
};

export const getStatusText = (status: string): string => {
  switch (status) {
    case 'scheduled':
      return 'Scheduled';
    case 'in_progress':
      return 'In Progress';
    case 'completed':
      return 'Completed';
    case 'missed':
      return 'Missed';
    case 'cancelled':
      return 'Cancelled';
    default:
      return status;
  }
};

export const getInitials = (fullName: string): string => {
  return fullName
    .split(' ')
    .map(name => name.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2);
};
