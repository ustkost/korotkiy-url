<script>
let a = $state(0)

const onclick = () => {
  a++;
}
</script>

<div>
  <button onclick={onclick}>
  click me
  </button>
  {a}
</div>
