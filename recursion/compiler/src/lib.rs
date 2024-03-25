#![allow(clippy::type_complexity)]
#![allow(clippy::needless_range_loop)]
extern crate alloc;

pub mod asm;
pub mod gnark;
pub mod ir;
pub mod util;
pub mod verifier;

pub mod prelude {
    pub use crate::asm::AsmCompiler;
    pub use crate::ir::*;
}
