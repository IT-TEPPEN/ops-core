import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { VariableForm } from './VariableForm';
import { VariableDefinition } from '../../types/domain';

describe('VariableForm', () => {
  const mockOnChange = vi.fn();
  const mockOnValidate = vi.fn();

  const sampleVariables: VariableDefinition[] = [
    {
      name: 'server_name',
      label: 'Server Name',
      description: 'The target server name',
      type: 'string',
      required: true,
      default_value: 'localhost'
    },
    {
      name: 'port',
      label: 'Port',
      description: 'Server port number',
      type: 'number',
      required: true,
      default_value: 8080
    },
    {
      name: 'enable_ssl',
      label: 'Enable SSL',
      type: 'boolean',
      required: false,
      default_value: true
    }
  ];

  beforeEach(() => {
    mockOnChange.mockClear();
    mockOnValidate.mockClear();
  });

  it('renders all variables', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
      />
    );

    expect(screen.getByLabelText('Server Name')).toBeInTheDocument();
    expect(screen.getByLabelText('Port')).toBeInTheDocument();
    expect(screen.getByLabelText('Enable SSL')).toBeInTheDocument();
  });

  it('renders required indicator for required fields', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
      />
    );

    const requiredMarkers = screen.getAllByLabelText('required');
    expect(requiredMarkers).toHaveLength(2); // server_name and port are required
  });

  it('renders descriptions when provided', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
      />
    );

    expect(screen.getByText('The target server name')).toBeInTheDocument();
    expect(screen.getByText('Server port number')).toBeInTheDocument();
  });

  it('calls onChange when string input changes', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{ server_name: 'localhost' }}
        onChange={mockOnChange}
      />
    );

    const input = screen.getByLabelText('Server Name');
    fireEvent.change(input, { target: { value: 'prod-server' } });

    expect(mockOnChange).toHaveBeenCalledWith('server_name', 'prod-server');
  });

  it('calls onChange when number input changes', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{ port: 8080 }}
        onChange={mockOnChange}
      />
    );

    const input = screen.getByLabelText('Port');
    fireEvent.change(input, { target: { value: '3000' } });

    expect(mockOnChange).toHaveBeenCalledWith('port', 3000);
  });

  it('calls onChange when checkbox changes', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{ enable_ssl: true }}
        onChange={mockOnChange}
      />
    );

    const checkbox = screen.getByLabelText('Enable SSL');
    fireEvent.click(checkbox);

    expect(mockOnChange).toHaveBeenCalledWith('enable_ssl', false);
  });

  it('displays default values', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
      />
    );

    const serverNameInput = screen.getByLabelText('Server Name') as HTMLInputElement;
    expect(serverNameInput.value).toBe('localhost');

    const portInput = screen.getByLabelText('Port') as HTMLInputElement;
    expect(portInput.value).toBe('8080');

    const sslCheckbox = screen.getByLabelText('Enable SSL') as HTMLInputElement;
    expect(sslCheckbox.checked).toBe(true);
  });

  it('does not render validate button when onValidate is not provided', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
      />
    );

    expect(screen.queryByText('Validate')).not.toBeInTheDocument();
  });

  it('renders validate button when onValidate is provided', () => {
    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
        onValidate={mockOnValidate}
      />
    );

    expect(screen.getByText('Validate')).toBeInTheDocument();
  });

  it('calls onValidate when validate button is clicked', async () => {
    mockOnValidate.mockResolvedValue(true);

    render(
      <VariableForm
        variables={sampleVariables}
        values={{}}
        onChange={mockOnChange}
        onValidate={mockOnValidate}
      />
    );

    const validateButton = screen.getByText('Validate');
    fireEvent.click(validateButton);

    expect(mockOnValidate).toHaveBeenCalled();
  });

  it('renders nothing when variables array is empty', () => {
    const { container } = render(
      <VariableForm
        variables={[]}
        values={{}}
        onChange={mockOnChange}
      />
    );

    expect(container.firstChild).toBeNull();
  });

  it('renders date input for date type', () => {
    const dateVariable: VariableDefinition = {
      name: 'backup_date',
      label: 'Backup Date',
      type: 'date',
      required: true,
      default_value: '2024-01-01'
    };

    render(
      <VariableForm
        variables={[dateVariable]}
        values={{}}
        onChange={mockOnChange}
      />
    );

    const input = screen.getByLabelText('Backup Date') as HTMLInputElement;
    expect(input.type).toBe('date');
  });
});
