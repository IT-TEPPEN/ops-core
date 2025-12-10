import { describe, it, expect } from 'vitest';
import { substituteVariables } from './variableSubstitution';

describe('substituteVariables', () => {
  it('should replace single variable', () => {
    const content = 'Connect to {{server_name}}';
    const values = { server_name: 'prod-server' };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Connect to prod-server');
  });

  it('should replace multiple variables', () => {
    const content = 'Connect to {{server_name}} and backup to {{backup_path}}';
    const values = {
      server_name: 'prod-server',
      backup_path: '/backup/db'
    };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Connect to prod-server and backup to /backup/db');
  });

  it('should replace multiple occurrences of the same variable', () => {
    const content = 'Server: {{server_name}}, Connect to {{server_name}}';
    const values = { server_name: 'prod-server' };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Server: prod-server, Connect to prod-server');
  });

  it('should leave unreplaced variables as-is when value not provided', () => {
    const content = 'Connect to {{server_name}} and backup to {{backup_path}}';
    const values = { server_name: 'prod-server' };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Connect to prod-server and backup to {{backup_path}}');
  });

  it('should handle empty values object', () => {
    const content = 'Connect to {{server_name}}';
    const values = {};
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Connect to {{server_name}}');
  });

  it('should handle content without variables', () => {
    const content = 'This is plain text';
    const values = { server_name: 'prod-server' };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('This is plain text');
  });

  it('should convert number values to string', () => {
    const content = 'Retention: {{retention_days}} days';
    const values = { retention_days: 30 };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Retention: 30 days');
  });

  it('should convert boolean values to string', () => {
    const content = 'Compression enabled: {{enable_compression}}';
    const values = { enable_compression: true };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Compression enabled: true');
  });

  it('should handle null and undefined values', () => {
    const content = 'Value: {{test_value}}';
    const values = { test_value: null };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Value: ');
  });

  it('should handle undefined values', () => {
    const content = 'Value: {{test_value}}';
    const values = { test_value: undefined };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Value: ');
  });

  it('should handle missing variables from values object', () => {
    const content = 'Value: {{missing_var}} and {{present_var}}';
    const values = { present_var: 'hello' };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Value: {{missing_var}} and hello');
  });

  it('should escape special regex characters in variable names', () => {
    const content = 'Path: {{base.path}}';
    const values = { 'base.path': '/home/user' };
    
    const result = substituteVariables(content, values);
    
    expect(result).toBe('Path: /home/user');
  });
});
