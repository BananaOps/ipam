// REST API client service for IPAM backend communication

import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import type {
  Subnet,
  SubnetFilters,
  CreateSubnetRequest,
  UpdateSubnetRequest,
  SubnetListResponse,
  APIError,
} from '../types';
import { LocationType, CloudProviderType } from '../types';
import { API_BASE_URL } from '../config/constants';
import { translateError } from '../utils/errorMessages';

/**
 * Transform API response from snake_case to camelCase
 */
function transformSubnetFromAPI(apiSubnet: any): Subnet {
  return {
    id: apiSubnet.id,
    cidr: apiSubnet.cidr,
    name: apiSubnet.name,
    description: apiSubnet.description,
    location: apiSubnet.location,
    locationType: apiSubnet.location_type as LocationType,
    cloudInfo: apiSubnet.cloud_info ? {
      provider: apiSubnet.cloud_info.provider as CloudProviderType,
      region: apiSubnet.cloud_info.region,
      accountId: apiSubnet.cloud_info.account_id,
    } : undefined,
    details: {
      address: apiSubnet.details.address,
      netmask: apiSubnet.details.netmask,
      wildcard: apiSubnet.details.wildcard,
      network: apiSubnet.details.network,
      type: apiSubnet.details.type,
      broadcast: apiSubnet.details.broadcast,
      hostMin: apiSubnet.details.host_min,
      hostMax: apiSubnet.details.host_max,
      hostsPerNet: apiSubnet.details.hosts_per_net,
      isPublic: apiSubnet.details.is_public,
    },
    utilization: {
      totalIps: apiSubnet.utilization.total_ips,
      allocatedIps: apiSubnet.utilization.allocated_ips,
      utilizationPercent: apiSubnet.utilization.utilization_percent,
    },
    createdAt: apiSubnet.created_at,
    updatedAt: apiSubnet.updated_at,
  };
}

/**
 * APIClient class for communicating with the IPAM backend REST API
 * Handles all HTTP operations with proper error handling and interceptors
 */
class APIClient {
  private axiosInstance: AxiosInstance;
  private authToken: string | null = null;
  private maxRetries: number = 3;
  private retryDelay: number = 1000; // milliseconds

  constructor(baseURL: string = API_BASE_URL) {
    // Log the base URL for debugging
    console.log('[APIClient] Initializing with baseURL:', baseURL);
    
    // Create axios instance with default configuration
    this.axiosInstance = axios.create({
      baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Set up request interceptor for authentication
    this.axiosInstance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // Add authentication token if available (for future use)
        if (this.authToken && config.headers) {
          config.headers.Authorization = `Bearer ${this.authToken}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Set up response interceptor for error handling
    this.axiosInstance.interceptors.response.use(
      (response) => {
        // Return successful responses as-is
        return response;
      },
      (error: AxiosError) => {
        // Transform errors into consistent APIError format
        const apiError = this.handleError(error);
        return Promise.reject(apiError);
      }
    );
  }

  /**
   * Set authentication token for future requests
   * @param token - JWT or other authentication token
   */
  setAuthToken(token: string | null): void {
    this.authToken = token;
  }

  /**
   * Get current authentication token
   */
  getAuthToken(): string | null {
    return this.authToken;
  }

  /**
   * Configure retry settings
   * @param maxRetries - Maximum number of retry attempts
   * @param retryDelay - Delay between retries in milliseconds
   */
  setRetryConfig(maxRetries: number, retryDelay: number): void {
    this.maxRetries = maxRetries;
    this.retryDelay = retryDelay;
  }

  /**
   * Retry a failed request with exponential backoff
   * @param fn - Function to retry
   * @param retries - Number of retries remaining
   * @returns Promise with the result
   */
  private async retryRequest<T>(
    fn: () => Promise<T>,
    retries: number = this.maxRetries
  ): Promise<T> {
    try {
      return await fn();
    } catch (error) {
      if (retries <= 0) {
        throw error;
      }

      // Only retry on network errors or 5xx server errors
      const shouldRetry = this.shouldRetryError(error as AxiosError);
      if (!shouldRetry) {
        throw error;
      }

      // Calculate delay with exponential backoff
      const delay = this.retryDelay * Math.pow(2, this.maxRetries - retries);
      
      await new Promise((resolve) => setTimeout(resolve, delay));
      
      return this.retryRequest(fn, retries - 1);
    }
  }

  /**
   * Determine if an error should trigger a retry
   * @param error - Axios error object
   * @returns true if the request should be retried
   */
  private shouldRetryError(error: AxiosError): boolean {
    // Retry on network errors
    if (!error.response) {
      return true;
    }

    // Retry on 5xx server errors
    const status = error.response.status;
    if (status >= 500 && status < 600) {
      return true;
    }

    // Retry on 429 (Too Many Requests)
    if (status === 429) {
      return true;
    }

    // Don't retry on client errors (4xx)
    return false;
  }

  /**
   * Create a new subnet
   * @param data - Subnet creation data
   * @returns Created subnet with calculated properties
   */
  async createSubnet(data: CreateSubnetRequest): Promise<Subnet> {
    return this.retryRequest(async () => {
      // Transform camelCase to snake_case for backend API
      const requestData = {
        cidr: data.cidr,
        name: data.name,
        description: data.description,
        location: data.location,
        location_type: data.locationType,
        cloud_info: data.cloudInfo ? {
          provider: data.cloudInfo.provider,
          region: data.cloudInfo.region,
          account_id: data.cloudInfo.accountId,
        } : undefined,
      };
      
      console.log('[APIClient] Creating subnet with data:', requestData);
      console.log('[APIClient] Request URL:', this.axiosInstance.defaults.baseURL + '/subnets');
      
      const response = await this.axiosInstance.post<any>('/subnets', requestData);
      return transformSubnetFromAPI(response.data);
    });
  }

  /**
   * List subnets with optional filters
   * @param filters - Optional filters for location, cloud provider, and search
   * @returns List of subnets matching the filters
   */
  async listSubnets(filters: SubnetFilters = {}): Promise<SubnetListResponse> {
    return this.retryRequest(async () => {
      const params = new URLSearchParams();
      
      if (filters.location) {
        params.append('location', filters.location);
      }
      if (filters.cloudProvider) {
        params.append('cloud_provider', filters.cloudProvider);
      }
      if (filters.searchQuery) {
        params.append('search', filters.searchQuery);
      }

      const response = await this.axiosInstance.get<any>('/subnets', {
        params,
      });
      
      // Transform API response from snake_case to camelCase
      return {
        subnets: response.data.subnets.map(transformSubnetFromAPI),
        totalCount: response.data.total_count,
      };
    });
  }

  /**
   * Get a specific subnet by ID
   * @param id - Subnet ID
   * @returns Subnet details
   */
  async getSubnet(id: string): Promise<Subnet> {
    return this.retryRequest(async () => {
      const response = await this.axiosInstance.get<any>(`/subnets/${id}`);
      return transformSubnetFromAPI(response.data);
    });
  }

  /**
   * Update an existing subnet
   * @param id - Subnet ID
   * @param data - Updated subnet data
   * @returns Updated subnet with recalculated properties
   */
  async updateSubnet(id: string, data: UpdateSubnetRequest): Promise<Subnet> {
    return this.retryRequest(async () => {
      // Transform camelCase to snake_case for backend API
      const requestData = {
        cidr: data.cidr,
        name: data.name,
        description: data.description,
        location: data.location,
        location_type: data.locationType,
        cloud_info: data.cloudInfo ? {
          provider: data.cloudInfo.provider,
          region: data.cloudInfo.region,
          account_id: data.cloudInfo.accountId,
        } : undefined,
      };
      
      console.log('[APIClient] Updating subnet with data:', requestData);
      
      const response = await this.axiosInstance.put<any>(`/subnets/${id}`, requestData);
      return transformSubnetFromAPI(response.data);
    });
  }

  /**
   * Delete a subnet
   * @param id - Subnet ID
   */
  async deleteSubnet(id: string): Promise<void> {
    return this.retryRequest(async () => {
      await this.axiosInstance.delete(`/subnets/${id}`);
    });
  }

  /**
   * Handle axios errors and transform them into APIError format
   * @param error - Axios error object
   * @returns Structured API error
   */
  private handleError(error: AxiosError): APIError {
    // Use the error translation utility for user-friendly messages
    const translation = translateError(error);
    
    // Check if error response exists and has expected structure
    if (error.response?.data) {
      const data = error.response.data as any;
      
      // If backend returns structured error, use translated message
      if (data.error && data.error.code) {
        return {
          code: data.error.code,
          message: translation.message,
          details: {
            original: data.error.message,
            suggestion: translation.suggestion || '',
            title: translation.title
          },
          timestamp: data.error.timestamp || Date.now(),
        };
      }
      
      // Handle direct error responses
      if (data.code && data.message) {
        return {
          code: data.code,
          message: translation.message,
          details: {
            original: data.message,
            suggestion: translation.suggestion || '',
            title: translation.title
          },
          timestamp: data.timestamp || Date.now(),
        };
      }
    }

    // For all other errors, use translated messages
    return {
      code: error.code || 'UNKNOWN_ERROR',
      message: translation.message,
      details: {
        original: error.message,
        suggestion: translation.suggestion || '',
        title: translation.title
      },
      timestamp: Date.now(),
    };
  }
}

// Export singleton instance for use throughout the application
export const apiClient = new APIClient();

// Export class for testing purposes
export default APIClient;
