export interface Server {
  id: string;
  name: string;
  total_cpu: number;
  total_ram: number;
  total_storage: number;
  bandwidth: number;
  monthly_price: number;
  price_per_hour: number;
  used_cpu: number;
  used_ram: number;
  used_storage: number;
}

export interface LoadSnapshot {
  timestamp: string;
  used_cpu: number;
  used_ram: number;
  used_storage: number;
}

export interface LeaseRequest {
  user_id: string;
  server_id: string;
  cpu: number;
  ram: number;
  storage: number;
}

export interface ReleaseRequest {
  user_id: string;
  server_id: string;
}