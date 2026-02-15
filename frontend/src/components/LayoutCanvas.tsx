"use client";

import { useRef, useState, useCallback } from "react";
import { Stage, Layer, Rect, Text, Transformer } from "react-konva";
import type { KonvaEventObject } from "konva/lib/Node";
import type Konva from "konva";
import type { LayoutItem } from "@/stores/useLayoutStore";
import { ITEM_TEMPLATES } from "./LayoutItemPalette";

type Props = {
  items: LayoutItem[];
  width: number;
  height: number;
  onItemsChange: (items: LayoutItem[]) => void;
};

function getColor(type: string): string {
  return ITEM_TEMPLATES.find((t) => t.type === type)?.color ?? "#888";
}

export function LayoutCanvas({ items, width, height, onItemsChange }: Props) {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const transformerRef = useRef<Konva.Transformer>(null);
  const stageRef = useRef<Konva.Stage>(null);

  const handleSelect = useCallback(
    (id: string) => {
      setSelectedId(id);
      // Transformer will attach via useEffect-like behavior in onTransformEnd
    },
    []
  );

  const handleStageClick = useCallback(
    (e: KonvaEventObject<MouseEvent | TouchEvent>) => {
      if (e.target === e.target.getStage()) {
        setSelectedId(null);
      }
    },
    []
  );

  const handleDragEnd = useCallback(
    (id: string, e: KonvaEventObject<DragEvent>) => {
      const updated = items.map((item) =>
        item.id === id
          ? { ...item, x: Math.round(e.target.x()), y: Math.round(e.target.y()) }
          : item
      );
      onItemsChange(updated);
    },
    [items, onItemsChange]
  );

  const handleTransformEnd = useCallback(
    (id: string, e: KonvaEventObject<Event>) => {
      const node = e.target as Konva.Rect;
      const scaleX = node.scaleX();
      const scaleY = node.scaleY();
      node.scaleX(1);
      node.scaleY(1);
      const updated = items.map((item) =>
        item.id === id
          ? {
              ...item,
              x: Math.round(node.x()),
              y: Math.round(node.y()),
              width: Math.round(node.width() * scaleX),
              height: Math.round(node.height() * scaleY),
              rotation: Math.round(node.rotation()),
            }
          : item
      );
      onItemsChange(updated);
    },
    [items, onItemsChange]
  );

  const handleDelete = useCallback(() => {
    if (!selectedId) return;
    onItemsChange(items.filter((item) => item.id !== selectedId));
    setSelectedId(null);
  }, [selectedId, items, onItemsChange]);

  // Attach transformer to selected node
  const attachTransformer = useCallback(
    (node: Konva.Rect | null) => {
      if (transformerRef.current) {
        if (node) {
          transformerRef.current.nodes([node]);
        } else {
          transformerRef.current.nodes([]);
        }
        transformerRef.current.getLayer()?.batchDraw();
      }
    },
    []
  );

  return (
    <div className="relative">
      {selectedId && (
        <button
          onClick={handleDelete}
          className="absolute top-2 right-2 z-10 rounded-lg bg-red-500 px-3 py-1 text-xs text-white hover:bg-red-600"
        >
          削除
        </button>
      )}
      <Stage
        ref={stageRef}
        width={width}
        height={height}
        onClick={handleStageClick}
        onTap={handleStageClick}
        className="rounded-lg border border-gray-200 bg-white"
      >
        <Layer>
          {/* Grid lines */}
          {Array.from({ length: Math.floor(width / 50) + 1 }).map((_, i) => (
            <Rect
              key={`vgrid-${i}`}
              x={i * 50}
              y={0}
              width={1}
              height={height}
              fill="#f0f0f0"
              listening={false}
            />
          ))}
          {Array.from({ length: Math.floor(height / 50) + 1 }).map((_, i) => (
            <Rect
              key={`hgrid-${i}`}
              x={0}
              y={i * 50}
              width={width}
              height={1}
              fill="#f0f0f0"
              listening={false}
            />
          ))}

          {items.map((item) => {
            const isSelected = item.id === selectedId;
            const color = getColor(item.type);
            return (
              <Rect
                key={item.id}
                id={item.id}
                x={item.x}
                y={item.y}
                width={item.width}
                height={item.height}
                rotation={item.rotation}
                fill={color}
                opacity={0.8}
                cornerRadius={4}
                draggable
                onClick={() => handleSelect(item.id)}
                onTap={() => handleSelect(item.id)}
                onDragEnd={(e) => handleDragEnd(item.id, e)}
                onTransformEnd={(e) => handleTransformEnd(item.id, e)}
                stroke={isSelected ? "#333" : undefined}
                strokeWidth={isSelected ? 2 : 0}
                ref={(node) => {
                  if (isSelected && node) {
                    attachTransformer(node);
                  }
                }}
              />
            );
          })}

          {/* Labels rendered on top */}
          {items.map((item) => (
            <Text
              key={`label-${item.id}`}
              x={item.x}
              y={item.y}
              width={item.width}
              height={item.height}
              text={item.label}
              fontSize={11}
              fill="#fff"
              align="center"
              verticalAlign="middle"
              listening={false}
              rotation={item.rotation}
            />
          ))}

          <Transformer
            ref={transformerRef}
            rotateEnabled={true}
            boundBoxFunc={(oldBox, newBox) => {
              if (newBox.width < 20 || newBox.height < 20) return oldBox;
              return newBox;
            }}
          />
        </Layer>
      </Stage>
    </div>
  );
}
