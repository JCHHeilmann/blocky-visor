<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import Input from "$lib/components/ui/Input.svelte";
  import Select from "$lib/components/ui/Select.svelte";

  interface Props {
    onsubmit: (domain: string, type: string) => void;
    loading?: boolean;
  }

  let { onsubmit, loading = false }: Props = $props();

  let domain = $state("");
  let type = $state("A");

  const recordTypes = [
    { value: "A", label: "A" },
    { value: "AAAA", label: "AAAA" },
    { value: "CNAME", label: "CNAME" },
    { value: "MX", label: "MX" },
    { value: "TXT", label: "TXT" },
    { value: "NS", label: "NS" },
    { value: "SOA", label: "SOA" },
    { value: "SRV", label: "SRV" },
    { value: "PTR", label: "PTR" },
  ];

  function handleSubmit(e: Event) {
    e.preventDefault();
    if (domain.trim()) {
      onsubmit(domain.trim(), type);
    }
  }
</script>

<form
  onsubmit={handleSubmit}
  class="flex flex-col gap-3 sm:flex-row sm:items-end"
>
  <div class="flex-1">
    <Input bind:value={domain} label="Domain" placeholder="example.com" />
  </div>
  <div class="w-full sm:w-32">
    <Select bind:value={type} options={recordTypes} label="Type" />
  </div>
  <Button type="submit" disabled={!domain.trim()} {loading}>Query</Button>
</form>
