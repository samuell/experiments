<?php

// -------------------------------------------------------------------------------------
//  What to run, when this file is executed as a script
// -------------------------------------------------------------------------------------

function main() {
	$fooProgram = new FooExample();
	$fooProgram->run();
}

// -------------------------------------------------------------------------------------
//  Program Examples
// -------------------------------------------------------------------------------------

class FooExample implements IProgram {
	function run() {
		// Initialize processes
		$foo_writer = new FooWriter();
		$foo_to_bar = new Foo2Bar();
		$printer = new Printer();

		// Connect network
		$foo_to_bar->in_foo = $foo_writer->out_foo;
		$printer->in_str = $foo_to_bar->out_bar;

		$net = new Network();
		$net->add_processes([$foo_writer, $foo_to_bar, $printer]);
		$net->run();
	}
}

// -------------------------------------------------------------------------------------
//  Processes
// -------------------------------------------------------------------------------------

class FooWriter implements IProcess {
	public $out_foo = null;

	public function __construct() {
		$this->out_foo = new Chan();
	}

	public function execute() {
		$this->out_foo->send('foo');
		$this->out_foo->close();

		return true;
	}
}

class Foo2Bar implements IProcess {
	public $in_foo = null;
	public $out_bar = null;

	public function __construct() {
		$this->in_foo = new Chan();
		$this->out_bar= new Chan();
	}

	public function execute() {
		$instr = $this->in_foo->recv();
		$outstr = str_replace('foo', 'bar', $instr);
		$this->out_bar->send($outstr);

		if ( $this->in_foo->closed() ) {
			$this->out_bar->close();
			return true;
		}
		return false;
	}
}

class Printer implements IProcess {
	public $in_str = null;

	public function __construct() {
		$this->in_str = new Chan();
	}

	public function execute() {
		$instr = $this->in_str->recv();
		echo "$instr\n";

		if ( $this->in_str->closed() ) {
			return true;
		}
		return false;
	}
}

// -------------------------------------------------------------------------------------
//  FBP Components
// -------------------------------------------------------------------------------------

class Network {
	protected $processes = [];

	public function add_processes($processes) {
		$this->processes = array_merge($this->processes, $processes);
	}

	public function run() {
		while ( count( $this->processes ) > 0 ) {
			foreach ( $this->processes as $i => $p ) {
				$done = $p->execute();
				if ( $done ) {
					unset($this->processes[$i]);
				}
			}
		}
	}
}

class Chan {
	protected $items = array();
	protected $closed = false;

	public function send($item) {
		array_unshift($this->items, $item);
	}

	public function recv() {
		return array_pop($this->items);
	}

	public function close() {
		$this->closed = true;
	}

	public function closed() {
		return $this->closed;
	}
}

interface INetwork {
	function run();
}

interface IProgram {
	function run();
}

interface IProcess {
	function execute();
}

// -------------------------------------------------------------------------------------
//  Run this file as a script
// -------------------------------------------------------------------------------------

main();