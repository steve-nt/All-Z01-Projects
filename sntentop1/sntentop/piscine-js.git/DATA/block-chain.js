function blockChain(data, prev) {
  const previousBlock = prev || { index: 0, hash: '0' }
  
  const block = {
    index: previousBlock.index + 1,
    data: data,
    prev: previousBlock,
    hash: '',
    chain: function(newData) {
      return blockChain(newData, block)
    }
  }
  
  block.hash = hashCode(block.index + previousBlock.hash + JSON.stringify(data))
  
  return block
}
