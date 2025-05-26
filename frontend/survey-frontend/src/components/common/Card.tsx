export default function Card({ title, children }: { title: string; children: React.ReactNode }) {
    return (
        <div className="p-4 border rounded shadow">
            <h2 className="text-xl font-bold">{title}</h2>
            {children}
        </div>
    );
}
