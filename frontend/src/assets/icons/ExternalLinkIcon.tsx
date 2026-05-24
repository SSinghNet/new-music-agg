interface IconProps {
  className?: string;
}

export default function ExternalLinkIcon({ className = "size-4" }: IconProps) {
  return (
    <svg className={className} viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden="true">
      <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 1.5H15.5M15.5 1.5V6.5M15.5 1.5L8 9M8.5 4H4a1 1 0 00-1 1v11a1 1 0 001 1h11a1 1 0 001-1v-4.5" />
    </svg>
  );
}
