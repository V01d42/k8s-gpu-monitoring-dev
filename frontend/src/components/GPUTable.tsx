import { useMemo } from 'react';
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  flexRender,
  createColumnHelper,
  type ColumnDef,
} from '@tanstack/react-table';
import { ArrowUpDown } from 'lucide-react';

import { GPUMetrics } from '@/types/gpu';
import { 
  formatPercentage, 
  formatBytes, 
  formatTemperature, 
  getUtilizationColor, 
  getTemperatureColor,
  cn 
} from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface GPUTableProps {
  data: GPUMetrics[];
  isLoading?: boolean;
  error?: Error | null;
}

const columnHelper = createColumnHelper<GPUMetrics>();

export function GPUTable({ data, isLoading, error }: GPUTableProps) {
  const columns = useMemo<ColumnDef<GPUMetrics, any>[]>(() => [
    columnHelper.accessor('node_name', {
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="h-auto p-0 font-semibold"
        >
          ノード名
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => (
        <div className="font-medium">{row.getValue('node_name')}</div>
      ),
    }),
    
    columnHelper.accessor('gpu_index', {
      header: 'GPU#',
      cell: ({ row }) => (
        <div className="text-center font-mono">
          {row.getValue('gpu_index')}
        </div>
      ),
    }),
    
    columnHelper.accessor('utilization', {
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="h-auto p-0 font-semibold"
        >
          GPU利用率
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const utilization = row.getValue('utilization') as number;
        return (
          <div className="flex items-center space-x-2">
            <div className="w-16 bg-muted rounded-full h-2">
              <div
                className={cn(
                  "h-2 rounded-full transition-all",
                  utilization < 30 ? "bg-green-500" :
                  utilization < 70 ? "bg-yellow-500" : "bg-red-500"
                )}
                style={{ width: `${Math.min(utilization, 100)}%` }}
              />
            </div>
            <span className={cn("font-mono text-sm", getUtilizationColor(utilization))}>
              {formatPercentage(utilization)}
            </span>
          </div>
        );
      },
    }),

    columnHelper.accessor('memory_utilization', {
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="h-auto p-0 font-semibold"
        >
          メモリ利用率
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const memoryUtil = row.getValue('memory_utilization') as number;
        return (
          <div className="flex items-center space-x-2">
            <div className="w-16 bg-muted rounded-full h-2">
              <div
                className={cn(
                  "h-2 rounded-full transition-all",
                  memoryUtil < 50 ? "bg-green-500" :
                  memoryUtil < 80 ? "bg-yellow-500" : "bg-red-500"
                )}
                style={{ width: `${Math.min(memoryUtil, 100)}%` }}
              />
            </div>
            <span className={cn("font-mono text-sm", getUtilizationColor(memoryUtil))}>
              {formatPercentage(memoryUtil)}
            </span>
          </div>
        );
      },
    }),

    columnHelper.accessor('memory_used', {
      header: 'メモリ使用量',
      cell: ({ row }) => {
        const memoryUsed = row.getValue('memory_used') as number;
        const memoryTotal = row.original.memory_total;
        return (
          <div className="text-sm">
            <div className="font-mono">{formatBytes(memoryUsed * 1024 * 1024 * 1024)}</div>
            <div className="text-muted-foreground text-xs">
              / {formatBytes(memoryTotal * 1024 * 1024 * 1024)}
            </div>
          </div>
        );
      },
    }),

    columnHelper.accessor('temperature', {
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="h-auto p-0 font-semibold"
        >
          温度
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const temperature = row.getValue('temperature') as number;
        return (
          <span className={cn("font-mono", getTemperatureColor(temperature))}>
            {formatTemperature(temperature)}
          </span>
        );
      },
    }),
  ], []);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    initialState: {
      pagination: {
        pageSize: 10,
      },
    },
  });

  if (error) {
    return (
      <Card className="w-full">
        <CardContent className="p-6">
          <div className="text-center text-destructive">
            <p className="text-lg font-semibold">エラーが発生しました</p>
            <p className="text-sm mt-2">{error.message}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>GPU監視ダッシュボード</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="rounded-md border">
          <table className="w-full">
            <thead>
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id} className="border-b">
                  {headerGroup.headers.map((header) => (
                    <th
                      key={header.id}
                      className="h-12 px-4 text-left align-middle font-medium text-muted-foreground"
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(header.column.columnDef.header, header.getContext())}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody>
              {isLoading ? (
                <tr>
                  <td colSpan={columns.length} className="h-24 text-center">
                    <div className="flex items-center justify-center space-x-2">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
                      <span>読み込み中...</span>
                    </div>
                  </td>
                </tr>
              ) : table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <tr
                    key={row.id}
                    className="border-b transition-colors hover:bg-muted/50"
                  >
                    {row.getVisibleCells().map((cell) => (
                      <td key={cell.id} className="p-4 align-middle">
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </td>
                    ))}
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={columns.length} className="h-24 text-center">
                    データがありません
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
} 