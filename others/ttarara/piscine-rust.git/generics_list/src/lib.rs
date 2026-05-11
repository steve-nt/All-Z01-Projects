#[derive(Clone, Debug)]
pub struct List<T> {
    pub head: Option<Node<T>>,
}

#[derive(Clone, Debug)]
pub struct Node<T> {
    pub value: T,
    pub next: Option<Box<Node<T>>>,
}

impl<T> List<T> {
    pub fn new() -> Self {
        List { head: None }
    }

    pub fn push(&mut self, value: T) {
        let next_head = self.head.take().map(Box::new);
        self.head = Some(Node { value, next: next_head });
    }

    pub fn pop(&mut self) {
        if let Some(head) = self.head.take() {
            self.head = head.next.map(|boxed| *boxed);
        }
    }

    pub fn len(&self) -> usize {
        let mut count = 0;
        let mut current = self.head.as_ref();
        while let Some(node) = current {
            count += 1;
            current = node.next.as_ref().map(|b| b.as_ref());
        }
        count
    }
}

