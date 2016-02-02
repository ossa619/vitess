<?php
// DO NOT EDIT! Generated by Protobuf-PHP protoc plugin 1.0
// Source: tabletmanagerdata.proto
//   Date: 2016-01-22 01:34:35

namespace Vitess\Proto\Tabletmanagerdata {

  class SlaveWasRestartedRequest extends \DrSlump\Protobuf\Message {

    /**  @var \Vitess\Proto\Topodata\TabletAlias */
    public $parent = null;
    

    /** @var \Closure[] */
    protected static $__extensions = array();

    public static function descriptor()
    {
      $descriptor = new \DrSlump\Protobuf\Descriptor(__CLASS__, 'tabletmanagerdata.SlaveWasRestartedRequest');

      // OPTIONAL MESSAGE parent = 1
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 1;
      $f->name      = "parent";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Topodata\TabletAlias';
      $descriptor->addField($f);

      foreach (self::$__extensions as $cb) {
        $descriptor->addField($cb(), true);
      }

      return $descriptor;
    }

    /**
     * Check if <parent> has a value
     *
     * @return boolean
     */
    public function hasParent(){
      return $this->_has(1);
    }
    
    /**
     * Clear <parent> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\SlaveWasRestartedRequest
     */
    public function clearParent(){
      return $this->_clear(1);
    }
    
    /**
     * Get <parent> value
     *
     * @return \Vitess\Proto\Topodata\TabletAlias
     */
    public function getParent(){
      return $this->_get(1);
    }
    
    /**
     * Set <parent> value
     *
     * @param \Vitess\Proto\Topodata\TabletAlias $value
     * @return \Vitess\Proto\Tabletmanagerdata\SlaveWasRestartedRequest
     */
    public function setParent(\Vitess\Proto\Topodata\TabletAlias $value){
      return $this->_set(1, $value);
    }
  }
}
