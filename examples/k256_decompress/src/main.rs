#![no_main]

extern crate succinct_zkvm;
succinct_zkvm::entrypoint!(main);

extern "C" {
    fn syscall_k256_decompress(compressed_key: &mut [u8; 64], is_odd: bool);
}

pub fn main() {
    let mut compressed_key: [u8; 33] = [0; 33];
    succinct_zkvm::io::read_slice(&mut compressed_key);

    let mut decompressed_key: [u8; 64] = [0; 64];
    decompressed_key[..32].copy_from_slice(&compressed_key[1..]);
    let is_odd = match compressed_key[0] {
        2 => false,
        3 => true,
        _ => panic!("Invalid compressed key"),
    };
    unsafe {
        syscall_k256_decompress(&mut decompressed_key, is_odd);
    }

    let mut result: [u8; 65] = [0; 65];
    result[0] = 4;
    result[1..].copy_from_slice(&decompressed_key);

    succinct_zkvm::io::write_slice(&result);
}
