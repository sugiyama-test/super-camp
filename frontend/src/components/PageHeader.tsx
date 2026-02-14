type Props = {
  title: string;
  description?: string;
};

export function PageHeader({ title, description }: Props) {
  return (
    <header>
      <h1 className="text-2xl font-bold text-[var(--camp-green)]">{title}</h1>
      {description && (
        <p className="mt-1 text-sm text-gray-600">{description}</p>
      )}
    </header>
  );
}
